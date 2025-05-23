-- Generated by dbtpl for the django schema.

-- table auth_group
CREATE TABLE auth_group (
  id INT(11) AUTO_INCREMENT,
  name VARCHAR(150) NOT NULL,
  UNIQUE (name),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- table django_content_type
CREATE TABLE django_content_type (
  id INT(11) AUTO_INCREMENT,
  app_label VARCHAR(100) NOT NULL,
  model VARCHAR(100) NOT NULL,
  UNIQUE (app_label, model),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- table auth_permission
CREATE TABLE auth_permission (
  id INT(11) AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  content_type_id INT(11) NOT NULL REFERENCES django_content_type (id),
  codename VARCHAR(100) NOT NULL,
  UNIQUE (content_type_id, codename),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- table auth_group_permissions
CREATE TABLE auth_group_permissions (
  id BIGINT(20) AUTO_INCREMENT,
  group_id INT(11) NOT NULL REFERENCES auth_group (id),
  permission_id INT(11) NOT NULL REFERENCES auth_permission (id),
  UNIQUE (group_id, permission_id),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- index auth_group_permissio_permission_id_84c5c92e_fk_auth_perm
CREATE INDEX auth_group_permissio_permission_id_84c5c92e_fk_auth_perm ON auth_group_permissions (permission_id);

-- table auth_user
CREATE TABLE auth_user (
  id INT(11) AUTO_INCREMENT,
  password VARCHAR(128) NOT NULL,
  last_login DATETIME(6),
  is_superuser TINYINT(1) NOT NULL,
  username VARCHAR(150) NOT NULL,
  first_name VARCHAR(150) NOT NULL,
  last_name VARCHAR(150) NOT NULL,
  email VARCHAR(254) NOT NULL,
  is_staff TINYINT(1) NOT NULL,
  is_active TINYINT(1) NOT NULL,
  date_joined DATETIME(6) NOT NULL,
  UNIQUE (username),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- table auth_user_groups
CREATE TABLE auth_user_groups (
  id BIGINT(20) AUTO_INCREMENT,
  user_id INT(11) NOT NULL REFERENCES auth_user (id),
  group_id INT(11) NOT NULL REFERENCES auth_group (id),
  UNIQUE (user_id, group_id),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- index auth_user_groups_group_id_97559544_fk_auth_group_id
CREATE INDEX auth_user_groups_group_id_97559544_fk_auth_group_id ON auth_user_groups (group_id);

-- table auth_user_user_permissions
CREATE TABLE auth_user_user_permissions (
  id BIGINT(20) AUTO_INCREMENT,
  user_id INT(11) NOT NULL REFERENCES auth_user (id),
  permission_id INT(11) NOT NULL REFERENCES auth_permission (id),
  UNIQUE (user_id, permission_id),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- index auth_user_user_permi_permission_id_1fbb5f2c_fk_auth_perm
CREATE INDEX auth_user_user_permi_permission_id_1fbb5f2c_fk_auth_perm ON auth_user_user_permissions (permission_id);

-- table authors
CREATE TABLE authors (
  author_id BIGINT(20) AUTO_INCREMENT,
  name LONGTEXT NOT NULL,
  PRIMARY KEY (author_id)
) ENGINE=InnoDB;

-- table books
CREATE TABLE books (
  book_id BIGINT(20) AUTO_INCREMENT,
  isbn VARCHAR(255) NOT NULL,
  book_type INT(11) NOT NULL,
  title VARCHAR(255) NOT NULL,
  year INT(11) NOT NULL,
  available DATETIME(6) NOT NULL,
  books_author_id_fkey BIGINT(20) NOT NULL REFERENCES authors (author_id),
  PRIMARY KEY (book_id)
) ENGINE=InnoDB;

-- index books_books_author_id_fkey_73ac0c26_fk_authors_author_id
CREATE INDEX books_books_author_id_fkey_73ac0c26_fk_authors_author_id ON books (books_author_id_fkey);

-- table tags
CREATE TABLE tags (
  tag_id BIGINT(20) AUTO_INCREMENT,
  tag VARCHAR(50) NOT NULL,
  PRIMARY KEY (tag_id)
) ENGINE=InnoDB;

-- table books_tags
CREATE TABLE books_tags (
  id BIGINT(20) AUTO_INCREMENT,
  book_id BIGINT(20) NOT NULL REFERENCES books (book_id),
  tag_id BIGINT(20) NOT NULL REFERENCES tags (tag_id),
  UNIQUE (book_id, tag_id),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- index books_tags_tag_id_8d70b40a_fk_tags_tag_id
CREATE INDEX books_tags_tag_id_8d70b40a_fk_tags_tag_id ON books_tags (tag_id);

-- table django_admin_log
CREATE TABLE django_admin_log (
  id INT(11) AUTO_INCREMENT,
  action_time DATETIME(6) NOT NULL,
  object_id LONGTEXT,
  object_repr VARCHAR(200) NOT NULL,
  action_flag SMALLINT(5) UNSIGNED NOT NULL,
  change_message LONGTEXT NOT NULL,
  content_type_id INT(11) REFERENCES django_content_type (id),
  user_id INT(11) NOT NULL REFERENCES auth_user (id),
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- index django_admin_log_content_type_id_c4bce8eb_fk_django_co
CREATE INDEX django_admin_log_content_type_id_c4bce8eb_fk_django_co ON django_admin_log (content_type_id);

-- index django_admin_log_user_id_c564eba6_fk_auth_user_id
CREATE INDEX django_admin_log_user_id_c564eba6_fk_auth_user_id ON django_admin_log (user_id);

-- table django_migrations
CREATE TABLE django_migrations (
  id BIGINT(20) AUTO_INCREMENT,
  app VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  applied DATETIME(6) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- table django_session
CREATE TABLE django_session (
  session_key VARCHAR(40) NOT NULL,
  session_data LONGTEXT NOT NULL,
  expire_date DATETIME(6) NOT NULL,
  PRIMARY KEY (session_key)
) ENGINE=InnoDB;

-- index django_session_expire_date_a5c62663
CREATE INDEX django_session_expire_date_a5c62663 ON django_session (expire_date);
