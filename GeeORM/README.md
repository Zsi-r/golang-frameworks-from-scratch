# GeeORM

**GeeORM**: A [gorm](https://github.com/jinzhu/gorm)-like and [xorm](https://github.com/go-xorm/xorm)-like object relational mapping library

- 实现表的创建、删除、迁移和记录的增删查改
- 实现事务(transaction)的提交和回滚，并封装成接口，从而确保事务的原子性
- 记录的查询功能支持链式操作，以提高代码的简洁度和可读性
- 实现钩子机制，使框架可在增删查改前后自动触发用户的自定义方法，提高了框架的灵活性和扩展性

## 技术点总结

最主要的技术点：**go的反射机制**

#### Q1：如何实现钩子机制？

=> 钩子机制同样是通过反射来实现的，`s.RefTable().Model` 或 `value` 即当前会话正在操作的对象，使用 `MethodByName` 方法反射得到该对象的方法。若该方法是有效的，则使用`.Call(param)`调用该方法。

#### Q2：如何实现数据库迁移？

=>

- 使用两次自定义的`difference()` 来计算前后两个字段切片的差集。新表 - 旧表 = 新增字段，旧表 - 新表 = 删除字段。
- 使用 `ALTER` 语句新增字段。
- 使用创建新表并重命名的方式删除字段。

#### Q3：如何实现事务？

=>

- Go 语言标准库 database/sql 提供了支持事务的接口
- 之前直接使用 `sql.DB` 对象执行 SQL 语句，如果要支持事务，需要更改为 `sql.Tx` 执行。当 `tx` 不为空时，则使用 `tx` 执行 SQL 语句，否则使用 `db` 执行 SQL 语句。这样既兼容了原有的执行方式，又提供了对事务的支持。
- 封装事务的 Begin、Commit 和 Rollback 三个接口。
- 用户只需要将所有的操作放到一个回调函数中，作为入参传递给 `engine.Transaction()`，**发生任何错误，自动调用Rollback回滚，如果没有错误发生，则调用Commit提交**。