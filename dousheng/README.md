有这个语句，修改struct，数据库会自动添加列。表空时，就是建表

```
MySQLDB.AutoMigrate(&model.User{}, &model.Video{})
```

