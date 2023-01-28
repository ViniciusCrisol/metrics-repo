CREATE TABLE app_metrics(
    id         INT         NOT NULL AUTO_INCREMENT PRIMARY KEY,
    data       JSON        NOT NULL,
    app_name   VARCHAR(64) NOT NULL,
    created_at DATETIME    NOT NULL
);

CREATE INDEX app_metrics_ix ON app_metrics(app_name);
