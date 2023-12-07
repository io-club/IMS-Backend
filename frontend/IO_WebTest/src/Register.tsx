import { Link, useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import './border.css';


export default function Register() {
  const color="white";
  //TODO URL
  const testUrl = "http://101.200.63.44:40000";
  //数据存储
  const [formData, setFormData] = useState({
    type: 'outsiders',
    name: '',
    password: '',
    rePassword: '',
    email: '',
    verificationCode: "",
  });
  const [isCounting, setIsCounting] = useState(false);
  const [timer, setTimer] = useState(60);
  const [ok, setOk] = useState(false);
  const navigate = useNavigate();


  function handleCloseModal() {
    sendCodeToEmail(formData.email);
    console.log(formData.email);
    // setIsCounting(true);
    setTimer(60);
    setShowPopup(false);
    setSliderPosition(0); // 当关闭弹窗时重置进度条
    document.body.style.overflow = 'auto'; // 恢复页面滚动
    alert("验证成功!");
    setIsVerified(false);
    return null;
  }




  useEffect(() => {
    if (ok) {
      navigate('/login'); // 进行页面跳转
    }
  }, [ok, navigate]);


  //数据处理判断并设置
  const handleInputChange = (event:React.ChangeEvent<HTMLInputElement>, field:string) => {
    let value = event.target.value;
    //if用户输入数据判断，前端or后端
    setFormData({
      ...formData,
      [field]: value,
    });
  };
  //数据提交
  const handleSubmit = (event:React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const { type, name, password, rePassword, email, verificationCode } = formData;
    if (!name
      || !password
      || !rePassword
      || !email
      || password !== rePassword
      || !verificationCode) {
      alert('请填写完整信息,并确保您的密码输入一致!');
      return;
    }
    addUserInfoToBackend(type, name, email, verificationCode, password); //发送请求到后端
  };

  //提交数据到后端
  function addUserInfoToBackend(type:string, name:string, email:string, verificationCode:string,password:string) {
    const url = testUrl + '/user/Register'
    const formData = new FormData();
    formData.append('type', type);
    formData.append('name', name);
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
          return response.json();
        }
        throw new Error('Network response was not ok.');
      })
      .then((data) => {
        console.log('Success:', data);
        setOk(true);
        // 处理成功的响应
      })
      .catch((error) => {
        console.error('Error:', error);
        // 处理错误
      });
  }
  //获取邮箱验证码
  function sendCodeToEmail(email:string) {
    const url = testUrl + '/user/SendVerification'; // 后端地址
  
    const formData = new FormData();
    formData.append('email', email);
    formData.append('url', 'http://localhost:9999/login');
  
    let isTimedOut = false; // 标记超时状态
  
    const timeoutPromise = new Promise((_, reject) => {
      setTimeout(() => {
        if (!isTimedOut) {
          isTimedOut = true;
          reject(new Error('Request timed out'));
        }
      }, 10000); // 设置超时时间为 10 秒钟
    });
  
    const fetchPromise = fetch(url, {
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
        if (!isTimedOut) {
          console.log('Success:', data);
          // 从响应中提取所需数据
          const { id, type, name, email, status } = data.data;
          localStorage.setItem('userData', JSON.stringify({ id, type, name, email, status }));
        }
      })
      .catch((error) => {
        if (!isTimedOut) {
          console.error('Error:', error);
          // 处理错误
        }
      });
  
    // 等待 fetchPromise 和 timeoutPromise 任意一个完成
    return Promise.race([fetchPromise, timeoutPromise]);
  }
  

  const handleEmail = (e:React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    if (!formData.email) {
      alert('请输入邮箱！');
      return;
    } else {
      setShowPopup(true);
      document.body.style.overflow = 'hidden'; // 禁用页面滚动
    }
  };
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
  /*------------------*/

  const [isVerified, setIsVerified] = useState(false); // 状态用于验证是否完成
  const [showPopup, setShowPopup] = useState(false); // 控制弹窗显示与隐藏


  /*------------------*/
  const [sliderPosition, setSliderPosition] = useState(0);

  const handleSlide = (event: React.ChangeEvent<HTMLInputElement>) => {
    setIsVerified(false);
    setSliderPosition(Number(event.target.value)); // 需要将滑块的值转换为数字，以便进行比较
  
    // 设置验证逻辑，可以根据滑块的位置来判断是否验证成功
    // 这里设置的是滑块位置大于某个阈值时认为验证通过
    if (Number(event.target.value) > 99) {
      setIsVerified(true);
      setShowPopup(false);
    }
  };
  

  return (
    <div className="container" style={{height:'600px'}}>
      <span></span>
      {showPopup && <div className="modal">

        <div className="popup">
          {/* 这里放置弹窗内容 */}

          <div style={{ width: '300px', margin: '20px auto' }}>
            <input
              type="range"
              min="0"
              max="100"
              value={sliderPosition}
              onChange={(value) => handleSlide(value)}
              style={{ width: '100%' }}
            />
          </div>
        </div>
      </div>}
      {showPopup && <div className="overlay" />}
      {isVerified && handleCloseModal()}


      <form className="form" onSubmit={handleSubmit}>
        <h2 style={{color:color}}>注册</h2>
        <div className="form-control">
          <label htmlFor="username">用户名</label>
          <input
            type="text"
            id="username"
            placeholder="请输入用户名"
            value={formData.name}
            onChange={(e) => handleInputChange(e, 'name')} />
        </div>

        <div className="form-control">
          <label style={{color:color}} htmlFor="email">邮箱</label>
          <input
            type="email"
            id="email"
            placeholder="请输入邮箱"
            value={formData.email}
            onChange={(e) => handleInputChange(e, 'email')} />
        </div>
        <button
          className='sendCodeBtn'
          id='sendCodeBtn'
          onClick={(e) => handleEmail(e)}
          type='button'
          disabled={isCounting && timer > 0}>
          {isCounting && timer > 0 ? `倒计时 ${timer}s` : '发送验证码'}
        </button>


        <div className="form-control">
          <label style={{color:color}} htmlFor="email">验证码</label>
          <input
            type="verificationCode"
            id="verificationCode"
            placeholder="请输入邮箱验证码"
            value={formData.verificationCode}
            onChange={(e) => handleInputChange(e, 'verificationCode')} />
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

        <div className="form-control">
          <label style={{color:color}} htmlFor="confirmPassword">确认密码</label>
          <input
            type="password"
            id="confirmPassword"
            placeholder="请再次输入密码"
            value={formData.rePassword}
            onChange={(e) => handleInputChange(e, 'rePassword')} />
        </div>
        <button type="submit">注册</button>
        <p style={{color:color}}>已有账户？ <Link to="/">登录</Link></p>
      </form>
    </div>
  );
}
