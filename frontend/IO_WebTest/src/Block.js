import React, { useState } from 'react';

const SliderVerification = () => {
  const [isVerified, setIsVerified] = useState(false);
  const [sliderPosition, setSliderPosition] = useState(0);

  const handleSlide = (event) => {
    const currentPosition = event.target.value;
    setSliderPosition(currentPosition);

    // 设置验证逻辑，可以根据滑块的位置来判断是否验证成功
    // 这里设置的是滑块位置大于某个阈值时认为验证通过
    if (currentPosition > 70) {
      setIsVerified(true);
    } else {
      setIsVerified(false);
    }
  };

  const handleVerification = () => {
    // 在这里执行验证通过后的操作
    if (isVerified) {
      alert('验证通过，执行下一步操作！');
      // 在这里执行你的下一步操作
    } else {
      alert('请拖动滑块至正确位置！');
    }
  };

  return (
    <div>
      <h2>滑块验证</h2>
      <div style={{ width: '300px', margin: '20px auto' }}>
        <input
          type="range"
          min="0"
          max="100"
          value={sliderPosition}
          onChange={handleSlide}
          style={{ width: '100%' }}
        />
      </div>
      <button onClick={handleVerification}>验证</button>
    </div>
  );
};

export default SliderVerification;
