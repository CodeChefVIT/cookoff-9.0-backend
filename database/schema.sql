CREATE TABLE users (
	id UUID NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE, 
	reg_no TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	role TEXT NOT NULL,
	round_qualified INTEGER NOT NULL DEFAULT 0,
	score INTEGER DEFAULT 0,
	name TEXT NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE questions (
	id UUID NOT NULL UNIQUE,
	description TEXT,
	title TEXT,
	input_format TEXT,
	points INTEGER,
	round INTEGER NOT NULL,
	constraints TEXT,
	output_format TEXT,
	PRIMARY KEY(id)
);

CREATE TABLE submissions (
	id UUID NOT NULL UNIQUE,
	question_id UUID NOT NULL,
	testcases_passed INTEGER DEFAULT 0,
	testcases_failed INTEGER DEFAULT 0,
	runtime DECIMAL,
	submission_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	testcase_id UUID,
	language_id INTEGER NOT NULL,
	description TEXT,
	memory INTEGER,
	user_id UUID,
	status TEXT,
	PRIMARY KEY(id)
);

CREATE TABLE testcases (
	id UUID NOT NULL UNIQUE,
	expected_output TEXT NOT NULL,
	memory TEXT NOT NULL,
	input TEXT NOT NULL,
	hidden BOOLEAN NOT NULL,
	runtime DECIMAL NOT NULL,
	question_id UUID NOT NULL,
	PRIMARY KEY(id)
);

-- Foreign keys
ALTER TABLE submissions
ADD FOREIGN KEY(question_id) REFERENCES questions(id)
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE testcases
ADD FOREIGN KEY(question_id) REFERENCES questions(id)
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE submissions
ADD FOREIGN KEY(testcase_id) REFERENCES testcases(id)
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE submissions
ADD FOREIGN KEY(user_id) REFERENCES users(id)
ON UPDATE NO ACTION ON DELETE NO ACTION;
