package v1

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"

	v1 "github.com/m1ckswagger/super-duper-survey/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceServer struct {
	v1.UnimplementedUserServiceServer
	db *sql.DB
}

// NewUserServiceServer creates a new User service
func NewUserServiceServer(db *sql.DB) v1.UserServiceServer {
	return &userServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *userServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *userServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database -> "+err.Error())
	}
	return c, nil
}

// Create new todo task
func (s *userServiceServer) Register(ctx context.Context, req *v1.UserRegisterRequest) (*v1.UserRegisterResponse, error) {
	// check for valid API version requested
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	md5pass := md5.Sum([]byte(req.User.Password))
	// insert user into db
	res, err := c.ExecContext(ctx, "INSERT INTO Users(`Email`, `FirstName`, `LastName`, `Password`, `IsAdmin`, `IsSuperUser`) VALUES(?, ?, ?, ?, ?, ?)",
		req.User.Email, req.User.Firstname, req.User.Lastname, md5pass, req.User.IsAdmin, req.User.IsSuperuser)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Users-> "+err.Error())
	}

	// get ID of created User
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve ID for created user-> "+err.Error())
	}

	return &v1.UserRegisterResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

func (s *userServiceServer) View(ctx context.Context, req *v1.UserViewRequest) (*v1.UserViewResponse, error) {
	// check API version
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query User by ID
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Email`, `FirstName`, `LastName`, `Password`, `IsAdmin`, `IsSuperUser` FROM Users WHERE `ID`=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Users-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Users-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("User with ID='%d' cannot be found",
			req.Id))
	}

	// data for User
	var usr v1.User

	if err := rows.Scan(&usr.Id, &usr.Email, &usr.Firstname, &usr.Lastname, &usr.Password, &usr.IsAdmin, &usr.IsSuperuser); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from User row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple User rows with ID='%d'", req.Id))
	}

	return &v1.UserViewResponse{
		Api:  apiVersion,
		User: &usr,
	}, nil
}

func (s *userServiceServer) Update(ctx context.Context, req *v1.UserUpdateRequest) (*v1.UserUpdateResponse, error) {
	// check API version
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	md5pass := md5.Sum([]byte(req.User.Password))
	res, err := c.ExecContext(ctx, "UPDATE Users SET `Email`=?, `FirstName`=?, `LastName`=?, `Password`=?, `IsAdmin`=?, `IsSuperUser`=?",
		req.User.Email, req.User.Firstname, req.User.Lastname, md5pass, req.User.IsAdmin, req.User.IsSuperuser)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update User-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("User with ID='%d' was not found", req.User.Id))
	}
	return &v1.UserUpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil
}

func (s *userServiceServer) ReadAll(ctx context.Context, req *v1.UserViewAllRequest) (*v1.UserViewAllResponse, error) {
	// check API version
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get User list
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Email`, `FirstName`, `LastName`, `Password`, `IsAdmin`, `IsSuperUser` FROM Users")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Users-> "+err.Error())
	}
	defer rows.Close()

	list := []*v1.User{}
	for rows.Next() {
		usr := new(v1.User)
		if err := rows.Scan(&usr.Id, &usr.Email, &usr.Firstname, &usr.Lastname, &usr.Password, &usr.IsAdmin, &usr.IsSuperuser); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from User row-> "+err.Error())
		}
		list = append(list, usr)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Users-> "+err.Error())
	}

	return &v1.UserViewAllResponse{
		Api:   apiVersion,
		Users: list,
	}, nil
}
