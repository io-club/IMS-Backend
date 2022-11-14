'''
这个文件包含了 casbin 的基本使用方法
'''
import casbin

e = casbin.Enforcer("./model.conf", "./policy.csv")


def check(sub, obj, act):
    '''
    sub：想要访问资源的用户
    obj：想要访问的资源
    act：想要执行的操作
    '''
    if e.enforce(sub, obj, act):
        print("通过")
    else:
        print("拒绝")


check("admin", "data1", "read")
check("admin", "data1", "write")
check("user", "data1", "read")
check("user", "data1", "write")
