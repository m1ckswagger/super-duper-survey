USE surveys;
CREATE TABLE IF NOT EXISTS surveys (
    id int NOT NULL,
    name varchar(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS answers (
    id int NOT NULL AUTO_INCREMENT,
    surveyID int NOT NULL,
    session_id varchar(255) NOT NULL,
    num int NOT NULL,
    answer varchar(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (surveyID) REFERENCES surveys(id)
);

CREATE TABLE IF NOT EXISTS ocai (
  id int NOT NULL AUTO_INCREMENT,
  session_id varchar(255) NOT NULL,
  department varchar(255) NOT NULL,
  num int NOT NULL,
  a_now int NOT NULL,
  b_now int NOT NULL,
  c_now int NOT NULL,
  d_now int NOT NULL,
  a_preferred int NOT NULL,
  b_preferred int NOT NULL,
  c_preferred int NOT NULL,
  d_preferred int NOT NULL,
  PRIMARY KEY (id)
);
