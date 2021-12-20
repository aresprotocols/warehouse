package sql

import (
	"github.com/jmoiron/sqlx"
)

var db_create = `
create database if not exists db_price;
`

var t_coin_history_info = `
create table if not exists t_coin_history_info
(
     id bigint(20) not NULL AUTO_INCREMENT primary key,
     symbol varchar(20) not null,
     timestamp integer not null,
     price decimal(28,8) not null,
     price_origin varchar(20) not null,
	 weight integer not null
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

var t_log_info = `
create table if not exists t_log_info
(
     id bigint(20) not NULL AUTO_INCREMENT primary key,
     client_ip varchar(20) not null,
	 request_time varchar(64) not null,
	 user_agent varchar(256) not null,
	 request_url varchar(512) not null,
	 response_time varchar(64) not null,
	 use_symbol integer not null,
	 request_response text not null,
     request_timestamp integer,
     response_timestamp integer
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

var t_http_error = `
create table if not exists t_http_error
(
     id bigint(20) not NULL AUTO_INCREMENT primary key,
     url varchar(128) not null,
	 symbol varchar(16) not null,
	 error varchar(1024) not null,
	 timestamp integer not null
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

var t_weight_info = `
create table if not exists t_weight_info
(
     id bigint(20) not NULL AUTO_INCREMENT primary key,
     symbol varchar(16) not null,
	 exchange varchar(16) not null,
	 weight integer not null
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

var t_update_price_history = `
create table if not exists t_update_price_history
(
    timestamp int         not null,
    symbol    varchar(20) not null,
    primary key (timestamp, symbol)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

func createOrderTables(db *sqlx.DB, dbName string) {
	db.MustExec(db_create)
	db.MustExec("USE " + dbName)
	db.MustExec(t_coin_history_info)
	db.MustExec(t_log_info)
	db.MustExec(t_http_error)
	db.MustExec(t_weight_info)
	db.MustExec(t_update_price_history)
}
