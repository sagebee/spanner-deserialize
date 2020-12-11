CREATE TABLE cards (
    id INT64 NOT NULL,  
    givenin INT64,      
    givento INT64,      
    timegiven INT64,    
    cardtype STRING(6)  
) PRIMARY KEY (id);

CREATE TABLE goals (
    id INT64 NOT NULL, 
    scoredin INT64,    
    scoredby INT64,    
    timescored INT64,  
    rating STRING(20)  
) PRIMARY KEY (id);

CREATE TABLE involves (
    match INT64 NOT NULL, 
    team INT64 NOT NULL   
) PRIMARY KEY (match, team);

CREATE TABLE matches (
    id INT64 NOT NULL,       
    city STRING(50) NOT NULL, 
    playedon DATE NOT NULL
) PRIMARY KEY (id);

CREATE TABLE players (
    id INT64 NOT NULL,
    name STRING(50) NOT NULL,
    birthday DATE,
    memberof INT64 NOT NULL,
    position STRING(20) 
) PRIMARY KEY (id);

CREATE TABLE teams (
    id INT64 NOT NULL,
    country STRING(50) NOT NULL
) PRIMARY KEY (id)
