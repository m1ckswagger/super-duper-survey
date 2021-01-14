package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	v1 "github.com/m1ckswagger/super-duper-survey/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type catalogServiceServer struct {
	db *sql.DB
}

// NewCatalogServiceServer creates a new Catalog service
func NewCatalogServiceServer(db *sql.DB) v1.CatalogServiceServer {
	return &catalogServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *catalogServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *catalogServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database -> "+err.Error())
	}
	return c, nil
}

// Create new todo task
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
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

	// get timestamp data and check if valid
	created, err := ptypes.Timestamp(req.Catalog.Created)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "field 'created' has invalid format-> "+err.Error())
	}

	// get timestamp data and check if valid
	updated, err := ptypes.Timestamp(req.Catalog.Updated)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "field 'updated' has invalid format-> "+err.Error())
	}

	// insert catalog into db
	res, err := c.ExecContext(ctx, "INSERT INTO Catalogs(`Title`, `Description`, `Created`, `Updated`) VALUES(?, ?, ?)",
		req.Catalog.Title, req.Catalog.Description, created, updated)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Catalogs-> "+err.Eror())
	}

	// get ID of created Catalog
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve ID for created catalog-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

func (s *catalogServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
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

	// query Catalog by ID
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Created`, `Updated` FROM Catalogs WHERE `ID`=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Catalogs-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Catalogs-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Catalog with ID='%d' cannot be found",
			req.Id))
	}

	// cd for Catalog data
	var cd v1.Catalog
	var created time.Time
	var updated time.Time

	if err := rows.Scan(&cd.Id, &cd.Title, &cd.Description, created, updated); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Catalog row-> "+err.Error())
	}

	cd.Created, err = ptypes.TimestampProto(created)
	if err != nil {
		return nil, status.Error(codes.Unknown, "field 'Created' has invalid format-> "+err.Error())
	}

	cd.Updated, err = ptypes.TimestampProto(updated)
	if err != nil {
		return nil, status.Error(codes.Unknown, "field 'Updated' has invalid format-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Catalog rows with ID='%d'", req.Id))
	}

	return &v1.ReadResponse{
		Api:     apiVersion,
		Catalog: &cd,
	}, nil
}

func (s *catalogServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
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

	/*
		// get timestamp data and check if valid
		created, err := ptypes.Timestamp(req.Catalog.Created)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "field 'created' has invalid format-> "+err.Error())
		}
	*/

	/*
		// get timestamp data and check if valid
		updated, err := ptypes.Timestamp(req.Catalog.Updated)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "field 'updated' has invalid format-> "+err.Error())
		}
	*/
	updated := time.Now()

	res, err := c.ExecContext(ctx, "UPDATE Catalogs SET `Title`=?, `Description`=?, `Updated`=?")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update Catalog-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Catalog with ID='%d' was not found", req.Catalog.Id))
	}
	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil
}

// Delete catalog
func (s *catalogServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
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

	res, err := c.ExecContext(ctx, "DELETE FROM Catalogs WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete Catalog-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Catalog with ID='%s' was not found", req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}

func (s *catalogServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
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

	// get Catalog list
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Created`, `Updated` FROM Catalogs")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Catalogs-> "+err.Error())
	}
	defer rows.Close()

	var created time.Time
	var updated time.Time

	list := []*v1.Catalog{}
	for rows.Next() {
		cata := new(v1.Catalog)
		if err := rows.Scan(&cata.Id, &cata.Title, &cata.Description, &created, &updated); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from Catalog row-> "+err.Error())
		}

		cata.Created, err = ptypes.TimestampProto(created)
		if err != nil {
			return nil, status.Error(codes.Unknown, "field 'Created' has invalid format-> "+err.Error())
		}
		cata.Updated, err = ptypes.TimestampProto(updated)
		if err != nil {
			return nil, status.Error(codes.Unknown, "field 'Updated' has invalid format-> "+err.Error())
		}
		list = append(list, cata)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Catalogs-> "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:      apiVersion,
		Catalogs: list,
	}, nil
}
