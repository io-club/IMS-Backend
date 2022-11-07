import casbin

e = casbin.Enforcer("./model.conf", "./policy.csv")

sub = "alice"  # 想要访问资源的用户
obj = "data1"  # 将要被访问的资源
act = "read"  # 用户对资源进行的操作


def check(sub, obj, act):
    if e.enforce(sub, obj, act):
        print("通过")
    else:
        print("拒绝")


check("admin", "data1", "read")
check("admin", "data1", "write")
check("user", "data1", "read")
check("user", "data1", "write")
