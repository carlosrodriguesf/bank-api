create function uuid()
    returns varchar(36)
    language plpgsql
as
$$
declare
    id varchar(36);
begin
    SELECT uuid_in(overlay(overlay(md5(random()::text || ':' || clock_timestamp()::text) placing '4' from 13) placing
                           to_hex(floor(random() * (11 - 8 + 1) + 8)::int)::text from 17)::cstring)
    into id;
    return id;
end;
$$;