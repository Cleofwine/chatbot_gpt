create table chat_records(
    id bigint not null AUTO_INCREMENT,
    account varchar(255) not null default '',
    group_id varchar(255) not null default '',
    user_msg text null,
    user_msg_tokens int not null default '0',
    user_msg_keywords varchar(1024) not null default '',
    ai_msg text null,
    ai_msg_tokens int not null default '0',
    req_tokens int not null default '0',
    create_at bigint not null default '0',
    enterprise_id varchar(255) not null default '',
    endpoint int not null default '0',
    endpoint_account varchar(255) not null default '',
    PRIMARY KEY (id),
    KEY index_create_at (create_at DESC)
) engine=InnoDB default charset=utf8mb4;