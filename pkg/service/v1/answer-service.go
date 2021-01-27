package v1

import (
	"context"
	"database/sql"
	"fmt"

	v1 "github.com/m1ckswagger/super-duper-survey/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type answerServiceServer struct {
	v1.UnimplementedAnswerServiceServer
	db *sql.DB
}

// NewAnswerServiceServer creates a new Answer service
func NewAnswerServiceServer(db *sql.DB) v1.AnswerServiceServer {
	return &answerServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *answerServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *answerServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database -> "+err.Error())
	}
	return c, nil
}

// Create new todo task
func (s *answerServiceServer) Create(ctx context.Context, req *v1.AnswerCreateRequest) (*v1.AnswerCreateResponse, error) {
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
	// insert answer into db
	res, err := c.ExecContext(ctx, "INSERT INTO Answers(`CatalogID`, `QuestionNum`, `OptionNum`, `SessionID`) VALUES(?, ?, ?, ?)",
		req.Answer.CatalogId, req.Answer.QuestionNum, req.Answer.OptionNum, req.Answer.SessionId)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Answers-> "+err.Error())
	}

	// get ID of created Answer
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve ID for created answer-> "+err.Error())
	}

	return &v1.AnswerCreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

func (s *answerServiceServer) View(ctx context.Context, req *v1.AnswerViewRequest) (*v1.AnswerViewResponse, error) {
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

	// query Answer by ID
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `CatalogID`, `QuestionID`, `OptionID`, `SessionID` FROM Answers WHERE `ID`=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Answers-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Answers-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Answer with ID='%d' cannot be found",
			req.Id))
	}

	// data for Answer
	var ans v1.Answer

	if err := rows.Scan(&ans.Id, &ans.CatalogId, &ans.QuestionNum, &ans.OptionNum); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Answer row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Answer rows with ID='%d'", req.Id))
	}

	return &v1.AnswerViewResponse{
		Api:    apiVersion,
		Answer: &ans,
	}, nil
}

func (s *answerServiceServer) ViewAll(ctx context.Context, req *v1.AnswerViewAllRequest) (*v1.AnswerViewAllResponse, error) {
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

	// get Answer list
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `CatalogID`, `QuestionNum`, `OptionNum`, `SessionID` FROM Answers")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Answers-> "+err.Error())
	}
	defer rows.Close()

	list := []*v1.Answer{}
	for rows.Next() {
		ans := new(v1.Answer)
		if err := rows.Scan(&ans.Id, &ans.CatalogId, &ans.QuestionNum, &ans.OptionNum, &ans.SessionId); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from Answer row-> "+err.Error())
		}
		list = append(list, ans)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Answers-> "+err.Error())
	}

	return &v1.AnswerViewAllResponse{
		Api:     apiVersion,
		Answers: list,
	}, nil
}
