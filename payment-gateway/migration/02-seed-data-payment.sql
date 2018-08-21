INSERT INTO payment_method (name) VALUES ('Bank Transfer');
INSERT INTO payment_method (name) VALUES ('Credit');
INSERT INTO payment_method (name) VALUES ('Paypal');

INSERT INTO payment_gateway (name) VALUES ('VietinBank');
INSERT INTO payment_gateway (name) VALUES ('Vietcombank');

INSERT INTO rule (gateway_id, regex) VALUES (2, '^(?P<bankname>Vietcombank):SD TK (?P<bank_number>\d{13})\ \+?\-?(?P<amount>\d{1,3}(\,?\d{3})*)(?P<currency>VND) luc (?P<day>\d{2})\-(?P<month>\d{2})\-(?P<year>\d{4}) (?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2}). SD (?P<balance>\d{1,3}(\,?\d{3})*)VND..*CODE(?P<transaction_code>\w{6,6}).*$');
 
INSERT INTO rule (gateway_id, regex) VALUES (1, '^VietinBank:(?P<day>\d{2})\/(?P<month>\d{2})\/(?P<year>\d{4}|\d{2})\s(?P<hour>\d{2}):(?P<minute>\d{2})\|TK:(?P<bank_number>\d{12})\|GD:\-?\+?(?P<amount>\d{1,3}(\,?\d{3})*)(?P<currency>VND)\|SDC:(?P<balance>\d{1,3}(\,?\d{3})*)VND\|ND:.*CODE(?P<transaction_code>\w{6,6}).*$');
