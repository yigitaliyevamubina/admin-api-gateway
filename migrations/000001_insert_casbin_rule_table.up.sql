-- Insert rules for unauthorized users
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'unauthorized', '/v1/swagger/*', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'unauthorized', '/v1/register', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'unauthorized', '/v1/verify/{email}/{code}', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'unauthorized', '/v1/login', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'unauthorized', '/v1/auth/login', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'unauthorized', '/v1/swagger/*', 'GET');

-- Insert rules for users
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('g', 'user', 'unauthorized', '*');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'user', '/v1/user/update/{id}', 'PUT');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'user', '/v1/user/delete/{id}', 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'user', '/v1/user/{id}', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'user', '/v1/user/password', 'POST');

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'user', '/v1/user/refresh', 'POST');

-- Insert rules for admins
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('g', 'admin', 'user', '*');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'admin', '/v1/user/create', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'admin', '/v1/users/{page}/{limit}/{filter}', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'admin', '/v1/auth/admins/{page}/{limit}', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'admin', '/v1/auth/admin/{id}', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'admin', '/v1/auth/update', 'PUT');

-- Insert rules for superadmins
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('g', 'superadmin', 'admin', '*');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'superadmin', '/v1/auth/create', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'superadmin', '/v1/auth/delete', 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'superadmin', '/v1/rbac/roles', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'superadmin', '/v1/rbac/policies/{role}', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'superadmin', '/v1/rbac/add/policy', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'superadmin', '/v1/rbac/delete/policy', 'DELETE');
