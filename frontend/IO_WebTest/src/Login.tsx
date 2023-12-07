import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import "./styles.css"
import './border.css';
import warningg from './resources/warning.gif'

export default function Login() {
  const color='white';
  const [speed,setSpeed]=useState("0");
 //当 speed 改变时更新动画
 useEffect(() => {
  const formElement = document.querySelector('.form') as HTMLElement;
  if (formElement) {
    formElement.style.animation = `test ${speed}s linear infinite`;
  }
}, [speed]);


  //用户名登录
  const [formData, setFormData] = useState({
    name: '',
    password: '',
  });
  //邮箱验证码登录
  const [formData2, setFormData2] = useState({
    email: '',
    verificationCode: '',
  });
  const [isCounting, setIsCounting] = useState(false);
  const [userContent, setUserContent] = useState(true);
  const [emailContent, setEmailContent] = useState(true);
  const [timer, setTimer] = useState(60);
  const [loginModel, setLoginModel] = useState('username');
  const testUrl = "http://101.200.63.44:40000";
  /*----------------------------------------------------- */
  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>, field:string) => {
    let value = event.target.value;
    setUserContent(true);
    //if用户输入数据判断，前端or后端
    setFormData({
      ...formData,
      [field]: value,
    });
  }

  const handleInputChange2 = (event:React.ChangeEvent<HTMLInputElement>, field:string) => {
    let value2 = event.target.value;
    setEmailContent(true);
    //if用户输入数据判断，前端or后端
    setFormData2({
      ...formData2,
      [field]: value2,
    });
  }
  /*----------------------------------------------------- */
  const handleSubmit = (event:React.ChangeEvent<HTMLFormElement>) => {
    if (loginModel === 'username') {
      event.preventDefault();
      console.log(formData);
      const { name, password } = formData;
      if (!name
        || !password) {
        //TODO: 前端提示
        setUserContent(false);
        return;
      }
      verifyByBackend(name, password);//发送请求到后端
    }
    if (loginModel === 'email') {
      event.preventDefault();
      console.log(formData2);
      const { email, verificationCode } = formData2;
      if (!email
        || !verificationCode) {
        setEmailContent(false);
        return;
      }
      LoginByBackend(email, verificationCode);//发送请求到后端
    } else {
      return;
    }
  };


  function LoginByBackend(email:string, verificationCode:string) {
    const url = testUrl + '/user/EmailLogin';
    const formData2 = new FormData();
    formData2.append('email', email);
    formData2.append('verificationCode', verificationCode);

    fetch(url, {
      method: 'POST',
      body: formData2,
      // 设置适当的 headers
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        }
        throw new Error('Network response was not ok.');
      })
      .then((data) => {
        console.log('Success:', data);
        alert("登录成功！");
        // 从响应中提取所需数据
      })
      .catch((error) => {
        console.error('Error:', error);
        // 处理错误
      });
  }


  function verifyByBackend(name:string, password:string) {
    const url = testUrl + '/user/NameLogin';//后端地址

    const formData = new FormData();
    formData.append('name', name);
    formData.append('password', password);

    fetch(url, {
      method: 'POST',
      body: formData,
      // 设置适当的 headers
    })
      .then((response) => {
        if (response.ok) {
          const accessToken = response.headers.get('Access-Token');
          const refreshToken = response.headers.get('Refresh-Token');
          if(accessToken&&refreshToken){
          console.log(accessToken);
          console.log(refreshToken);
          // 存储数据
          localStorage.setItem('accessToken', accessToken);
          localStorage.setItem('refreshToken', refreshToken);
          }
          return response.json();
        }
        throw new Error('Network response was not ok.');
      })
      .then((data) => {
        console.log('Success:', data);
        alert("登录成功！");
        // 从响应中提取所需数据
      })
      .catch((error) => {
        console.error('Error:', error);
            // const containerElement = document.querySelector('.container');
    // containerElement.style.animation = `falling 1s ease-in-out infinite alternate`;
    setSpeed("0.1");

    // 获取 body 元素
    const body = document.querySelector('body');
    if(body){
    // 设置背景为 GIF 图片 URL
    body.style.backgroundImage = `url(${warningg})`;
    }
    
    // 创建一个新的 style 元素
    const styleElement = document.createElement('style');
    document.head.appendChild(styleElement);
    
    // 定义新的样式规则
    const newStyleRule = `.container span::before {
      content: '';
      position: absolute;
      inset: 5px;
      background-image: url(${warningg});
      background-attachment: fixed;
      background-size: cover;
      background-position: center;
      filter: blur(2px);
    }`;
    
    // 将新的样式规则添加到新创建的 style 元素中
    if (styleElement.sheet) {
      styleElement.sheet.insertRule(newStyleRule);
    }
    
      });
  }

  function EchangeModel() {
    console.log("切换模型");
    setLoginModel('email');
  }
  function UchangeModel() {
    setLoginModel('username');
  }
  function handleEmail() {
    // setIsCounting(true);
    console.log(formData2.email);
    sendCodeToEmail(formData2.email);
  }

  function sendCodeToEmail(email:string) {
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
        // 从响应中提取所需数据
        const { id, type, name, email, status } = data.data;
        localStorage.setItem("userData", JSON.stringify({ id, type, name, email, status }));
      })
      .catch((error) => {
        console.error('Error:', error);
        // 处理错误
      });
  }

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

  return (
    <div>
      {loginModel === 'username' ? (

        <div className="container">
          <span></span>
          <span></span>
          <span></span>
          {/* <button onClick={handleData}>查看数据是否接收成功</button> */}
          <form className="form" onSubmit={handleSubmit}>
          <div className="form-content">
            <h2 style={{color:color}}>用户名登录</h2>
            <div className="form-control">
              <label style={{color:color}} htmlFor="username">用户名</label>
              <input
                type="text"
                id="name"
                placeholder="请输入用户名"
                value={formData.name}
                onChange={(e) => handleInputChange(e, 'name')} />
            </div>

            <div className="form-control">
              <label style={{color:color}} htmlFor="password">密码</label>
              <input
                type="password"
                id="password"
                placeholder="请输入密码"
                value={formData.password}
                onChange={(e) => handleInputChange(e, 'password')} />
            </div>
            {!userContent && <p style={{color:color}}>请填写完整信息!</p>}

            <button type="submit">登录</button>
            <div className='button-container'>
              <p>
              <button
                type="button"
                onClick={EchangeModel}
                style={{
                  background: 'none',
                  border: 'none',
                  padding: '0',
                  fontSize: 'inherit',
                  color: 'inherit',
                  textDecoration: 'underline',
                  textAlign: 'left',
                  cursor: 'pointer',
                }}> <p style={{color:color}}>邮箱登录</p></button>
                </p>
              <Link to="/findPass"
                style={{
                  background: 'none',
                  border: 'none',
                  display: 'inline',
                  marginLeft:220,
                  padding: '0',
                  fontSize: 'inherit',
                  color: 'inherit',
                  textDecoration: 'underline',
                  textAlign: 'right',
                  cursor: 'pointer',
                }}><p style={{color:color}}>找回密码</p></Link>

            </div>
            <p style={{color:color}}>还没有账户？
              <Link to="/register">注册</Link>
            </p>
            </div>
          </form>
        </div>

      ) : (
        <div className="container">
            <span></span>
          <span></span>
          <span></span>
          {/* <button onClick={handleData}>查看数据是否接收成功</button> */}
          <form className="form" onSubmit={handleSubmit}>
            <h2 style={{color:color}}>邮箱登录</h2>
            <div className="form-control">
              <label style={{color:color}} htmlFor="username">邮箱</label>
              <input
                type="text"
                id="username"
                placeholder="请输入邮箱"
                value={formData2.email}
                onChange={(e) => handleInputChange2(e, 'email')} />
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
                value={formData2.verificationCode}
                onChange={(e) => handleInputChange2(e, 'verificationCode')}
                placeholder="请输入验证码" />
            </div>
            {!emailContent && <p style={{color:color}}>请填写完整信息!</p>}

            <button type="submit">登录</button>
            <button
              type="button"
              onClick={UchangeModel}
              style={{
                background: 'none',
                border: 'none',
                padding: '0',
                fontSize: 'inherit',
                color: 'inherit',
                textDecoration: 'underline',
                textAlign: 'left',
                cursor: 'pointer',
              }}><p style={{color:color}}>用户名登录</p></button>
            <p style={{color:color}}>还没有账户？
              <Link to="/register">注册</Link>
            </p>
          </form>
        </div>
      )}
    </div>
  );

}


