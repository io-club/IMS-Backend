import React, { useState,useEffect, ReactHTML } from 'react';
export default function FindPass() {
  const color="white";
    const testUrl = "http://101.200.63.44:40000";
    const [isCounting, setIsCounting] = useState(false);
    const [timer, setTimer] = useState(60);
    const [content,setContent]=useState(true);
    const [formData, setFormData] = useState({
        email: '',
        verificationCode: '',
        password: '',
        rePassword: '',
      });

    useEffect(() => {
      let intervalId: NodeJS.Timeout | undefined;
    
        if (isCounting && timer > 0) {
          intervalId = setInterval(() => {
            setTimer((prevTimer) => prevTimer - 1);
          }, 1000);
        }
    
        return () => {
          clearInterval(intervalId);
        };
      }, [isCounting, timer]);
    
      useEffect(() => {
        if (timer === 0) {
          setIsCounting(false);
        }
      }, [timer]);

      const handleEmail=()=>{
        sendCodeToEmail(formData.email);
      }

      const sendCodeToEmail=(email:string)=> {
        const url = testUrl + '/user/SendVerification';//后端地址
    
        const formData = new FormData();
        formData.append('email', email);
        formData.append('url', "http://localhost:9999/login");
    
        fetch(url, {
          method: 'POST',
          body: formData,
          // 设置适当的 headers
        })
          .then((response) => {
            if (response.ok) {
              setIsCounting(true);
              return response.json();
            }
            throw new Error('Network response was not ok.');
          })
          .then((data) => {
            console.log('Success:', data);
          })
          .catch((error) => {
            console.error('Error:', error);
            // 处理错误
          });
      }

      const handleInputChange = (event:React.ChangeEvent<HTMLInputElement>, field:string) => {
        let value = event.target.value;
        setFormData({
          ...formData,
          [field]: value,
        });
      }

      const handleSubmit = (event:React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        console.log(formData);
        const { email,verificationCode,password } = formData;
        if (!email
          || !password
          || !verificationCode
          || password !== formData.rePassword) {
          //TODO: 前端提示
          setContent(false);
          return;
        }
        FindPassByBackend(email,verificationCode,password);//发送请求到后端
      }
      function FindPassByBackend(email:string,verificationCode:string,password:string) {
        const url = testUrl + '/user/RetrievePassword';//后端地址
    
        const formData = new FormData();
        formData.append('email', email);
        formData.append('verificationCode', verificationCode);
        formData.append('password', password);
    
        fetch(url, {
          method: 'POST',
          body: formData,
          // 设置适当的 headers
        })
          .then((response) => {
            if (response.ok) {
                alert("密码修改成功");
              return response.json();
            }
            throw new Error('Network response was not ok.');
          })
          .then((data) => {
            console.log('Success:', data);
          })
          .catch((error) => {
            console.error('Error:', error);
            // 处理错误
          });
      }
    


    return (
        <div className="container">
          <span></span>
          <form style={{color:color}} className="form" onSubmit={handleSubmit}>
            <h2 style={{color:color}}>找回密码</h2>
            <div className="form-control">
              <label htmlFor="username">邮箱</label>
              <input
                type="text"
                id="username"
                placeholder="请输入邮箱"
                value={formData.email}
                onChange={(e) => handleInputChange(e, 'email')}
                 />
            </div>
            <button
              className='sendCodeBtn'
              id='sendCodeBtn'
              onClick={handleEmail}
              type='button'
              disabled={isCounting && timer > 0}>
              {isCounting && timer > 0 ? `倒计时 ${timer}s` : '发送验证码'}
            </button>

            <div className="form-control">
              <label style={{color:color}} htmlFor="password">验证码</label>
              <input
                type="verificationCode"
                id="verificationCode"
                value={formData.verificationCode}
                onChange={(e) => handleInputChange(e, 'verificationCode')}
                placeholder="请输入验证码" />
            </div>
    
            <div  className="form-control">
              <label htmlFor="password">密码</label>
              <input
                type="password"
                id="password"
                placeholder="请输入密码"
                value={formData.password}
                onChange={(e) => handleInputChange(e, 'password')} />
            </div>

            <div className="form-control">
          <label style={{color:color}} htmlFor="confirmPassword">确认密码</label>
          <input
            type="password"
            id="confirmPassword"
            placeholder="请再次输入密码"
            value={formData.rePassword}
            onChange={(e) => handleInputChange(e, 'rePassword')} />
        </div>
        {!content && <p>请填写完整信息!</p>}

    
            <button type="submit">登录</button>
          </form>
        </div>
      );
}