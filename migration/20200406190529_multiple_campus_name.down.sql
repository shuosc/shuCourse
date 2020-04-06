alter table campus
    alter column name type varchar(128)
        using (name[1]);