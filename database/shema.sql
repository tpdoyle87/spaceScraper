CREATE TABLE articles (
                      Id SERIAL PRIMARY KEY,
                      Title VARCHAR(255),
                      URL VARCHAR(255),
                      Company VARCHAR(255),
                      Date VARCHAR(100),
                      Hash VARCHAR(255) UNIQUE
);
