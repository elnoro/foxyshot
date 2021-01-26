create database foxyshot_screenshots;
use foxyshot_screenshots;
create table foxyshot_screenshots
(
    id          int auto_increment,
    path        varchar(255) not null,
    description text         null,
    constraint foxyshot_screenshots_pk
        primary key (id)
);
create unique index foxyshot_screenshots_path_uindex
    on foxyshot_screenshots (path);
