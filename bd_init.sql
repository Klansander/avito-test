drop table if exists contents;
drop table if exists banners;



create table banners
(
    id         serial primary key,
    is_active  boolean                                not null,
    content    jsonb                                  not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null

);

create table contents
(
    id         serial primary key,
    banner_id  int
        constraint rel_banner references banners on delete cascade,
    tag_id     int not null,
    feature_id int not null

);

create unique index tag_feature on contents (tag_id, banner_id);

create or replace function fn_banner_ins(i_is_active bool, i_tag_id int[], i_feature_id int, i_created_at timestamp,
                                         i_updated_at timestamp,
                                         i_content json, out v_id int, out o_res int, out o_mes text) returns record
    language plpgsql
as
$$
declare

i int;
begin
    o_res = 0;
    o_mes = '';
insert into public.banners(is_active, content, created_at, updated_at)
values (i_is_active, i_content, i_created_at, i_updated_at)
    returning id into v_id;

for i in 1 .. array_length(i_tag_id, 1)
        loop
            if (select count(*)
                from public.contents c
                where c.tag_id = i_tag_id[i]
                  and c.feature_id = i_feature_id) != 0
            then
                o_res = 1;
                o_mes = 'Баннер i_tag_id= ' || i_tag_id::text || ', i_feature_id= ' || i_feature_id::text ||
                        ' уже существует';
                return;
end if;
insert into public.contents(banner_id, tag_id, feature_id)
values (v_id, i_tag_id[i], i_feature_id);
end loop;


    RETURN ;

end ;
$$;

--alter function fn_banner_ins( boolean, int[], int, json, out int) owner to grandeas;

create or replace function fn_banner_get(i_tag_id int, i_feature_id int, i_is_admin bool, out o_json json,
                                         out o_res int, out o_mes text) returns record
    language plpgsql
as
$$
begin

    o_res = 0;
    o_mes = '';
    if (select count(*)
        from public.contents c
        where c.tag_id = i_tag_id
          and c.feature_id = i_feature_id) = 0
    then
        o_res = 1;
o_mes = 'Баннер для тега не найден';
        return;
end if;

select json_build_object(
               'content', cb.content
       ) as p1
from (select b.content
      from contents
               inner join banners b on b.id = contents.banner_id
      where tag_id = i_tag_id
        and feature_id = i_feature_id
        and (is_active = true or i_is_admin = true)) as cb
    into o_json;

return;

end;
$$;


--alter function fn_banner_get(int,int, out json) owner to grandeas;

create or replace function fn_banner_list(i_tag_id int, i_feature_id int, i_limit int, i_offset int,
                                          out o_json json) returns json
    language plpgsql
as
$$
begin

select json_agg(t1.p1)
from (select json_build_object(
                     'content', cb.content,
                     'tag_id', array_agg(cb.tag_id),
                     'feature_id', cb.feature_id,
                     'banner_id', cb.banner_id,
                     'is_active', cb.is_active,
                     'created_at', cb.created_at,
                     'updated_at', cb.updated_at
             ) as p1
      from (select *
            from contents
                     inner join banners b on b.id = contents.banner_id
            where (i_tag_id is null or tag_id = i_tag_id)
              and (i_feature_id is null or feature_id = i_feature_id)
            order by contents.id
                limit i_limit offset i_offset) as cb
      group by feature_id, content, banner_id, is_active, created_at, updated_at) as t1
    into o_json;

return;

end;
$$;

--alter function fn_banner_list(int,int,int,int, json, int, text) owner to grandeas;

create or replace function fn_banner_del(i_banner_id int, out o_res int, out o_mes text) returns record
    language plpgsql
as
$$
begin
    o_res = 0;
    o_mes = '';
    if (select count(*)
        from public.banners b
        where b.id = i_banner_id) = 0
    then
        o_res = 1;
o_mes = 'Баннер для тега не найден';
        return;
end if;
delete
from public.banners b
where b.id = i_banner_id;

o_mes = 'Баннер успешно удален';
    return;

end;
$$;

--alter function fn_banner_del(int,int,int,int, json, int, text) owner to grandeas;


-- UPDATE banners
-- SET tag_id = 3

create or replace function fn_banner_get_by_id(i_banner_id int,
                                               out o_json json, out o_res int, out o_mes text) returns record
    language plpgsql
as
$$
begin

    o_mes = '';
    o_res = 0;


select json_build_object(
               'content', cb.content,
               'tag_id', array_agg(cb.tag_id),
               'feature_id', cb.feature_id,
               'banner_id', cb.banner_id,
               'is_active', cb.is_active,
               'created_at', cb.created_at,
               'updated_at', cb.updated_at
       )
from (select *
      from contents
               inner join banners b on b.id = contents.banner_id
      where (banner_id = i_banner_id)) as cb
group by feature_id, content, banner_id, is_active, created_at, updated_at
    into o_json;

return;

end;
$$;

--alter function fn_banner_get_by_id(int, json, int, text) owner to grandeas;