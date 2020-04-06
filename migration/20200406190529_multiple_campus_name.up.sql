alter table campus
    alter column name type varchar(128)[]
        using (array [name]);

update campus
set name=array ['本部','宝山']
where '本部' = any (name);