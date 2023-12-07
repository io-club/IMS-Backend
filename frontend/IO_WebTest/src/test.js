import React, { useState } from "react";
import { Button, Modal } from "antd";
import "./styles.css"
 
const MixWay = (props) => {
  const [visiable, setVisiable] = useState(false);
 
  const onOk = () => {
    console.log("编写自己的onOk逻辑");
    closeModal();
  };
 
  const closeModal = () => {
    setVisiable(false);
  };
 
  return (
    <>
      <Button onClick={() => setVisiable(true)}>按钮+弹窗</Button>
      <Modal
        title="按钮+弹窗"
        className="test-modal"
        open={visiable}
        onOk={onOk}
        onCancel={closeModal}
        afterClose={closeModal}
      >
        <p>弹窗内容</p>
      </Modal>
    </>
  );
};
 
export default MixWay;