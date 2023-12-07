import React, { useState } from 'react';
import { Layout, Menu } from 'antd';
import { useNavigate, Link } from 'react-router-dom';
import "../styles.css"
import { useTranslation } from 'react-i18next';

import OverviewComponent from './OverviewComponent';
// import EnergyComponent from './EnergyComponent';
import MapComponent from './MapComponent';
import LogComponent from './LogComponent';
import HistoryComponent from './HistoryComponent';
import HomeComponent from './HomeComponent';
import MediaComponent from './MediaComponent';
import TodoComponent from './TodoComponent';
import VideoMedia from './VideoMedia';
import PictureMedia from './PictureMedia';
import SettingComponent from './SettingComponent';
/*----------------------------------------*/
import overview from '../pic/overview.png';
import energy from '../pic/energy.png';
import history from '../pic/history.png';
import home from '../pic/home.png';
import log from '../pic/log.png';
import map from '../pic/map.png';
import media from '../pic/media.png';
import todo from '../pic/todo.png';
import closeMenu from '../pic/closeMenu.png';
import openMenu from '../pic/openMenu.png';
import leftarrow from '../pic/leftarrow.png';
import setting from '../pic/setting.png';
/*----------------------------------------*/

const { Sider, Content } = Layout;

export default function EnergyComponent() {
  const [selectedMenu, setSelectedMenu] = useState("energy");
  const [isOpen, setIsOpen] = useState(true);
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();

  const handleMenuClick = (key:string) => {
    setSelectedMenu(key);
  };

  let componentToDisplay;
  if (selectedMenu === 'overview') {
    componentToDisplay = <OverviewComponent />;
  } else if (selectedMenu === 'energy') {
    // componentToDisplay = <EnergyComponent />;
  } else if (selectedMenu === 'map') {
    componentToDisplay = <MapComponent />;
  } else if (selectedMenu === 'log') {
    componentToDisplay = <LogComponent />;
  } else if (selectedMenu === 'history') {
    componentToDisplay = <HistoryComponent />;
  } else if (selectedMenu === 'home') {
    componentToDisplay = <HomeComponent />;
  } else if (selectedMenu === 'media') {
    componentToDisplay = <MediaComponent/>;
  } else if (selectedMenu === 'todo') {
    componentToDisplay = <TodoComponent />;
  } else if (selectedMenu === 'videoMedia') {
    componentToDisplay = <VideoMedia />
  } else if (selectedMenu === 'pictureMedia') {
    componentToDisplay = <PictureMedia />
  }else if(selectedMenu==='setting'){
    componentToDisplay = <SettingComponent />;
  }
  function handleOpenClick() {
    setIsOpen(true);
  }
  function handleCloseClick() {
    setIsOpen(false);
  }
  function handleArrowClick() {
    window.history.back();
  }

  return (
    <Layout>
      <Sider
        width={isOpen ? 200 : 80}
        style={{ background: '#fff' }}>
        <h2 style={{ display: 'flex', alignItems: 'center' }}>
          <div onClick={!isOpen ? handleOpenClick : handleCloseClick}>
            <img src={!isOpen ? openMenu : closeMenu} className='show-menu' alt="Icon" />
          </div>
          &nbsp;&nbsp;
          {isOpen && <span>{t('io')}&nbsp;{t("club")}</span>}
        </h2>
        <hr />
        <Menu selectedKeys={[selectedMenu]} mode="inline" onClick={(e) => handleMenuClick(e.key)}>
                    <Menu.Item className='menu-item' key="overview" icon={<img src={overview} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/overview">{t('overview')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="energy" icon={<img src={energy} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/energy">{t('energy')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="map" icon={<img src={map} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/map">{t('map')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="log" icon={<img src={log} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/log">{t('log')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="history" icon={<img src={history} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/history">{t('history')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="home" icon={<img src={home} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/home">{t('home')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="media" icon={<img src={media} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/media">{t('media')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="todo" icon={<img src={todo} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/todo">{t('todo')}</Link>
                    </Menu.Item>
                    <Menu.Item className='menu-item' key="setting" icon={<img src={setting} />}
                        style={{ marginBottom: '20px' }}>
                        <Link to="/setting">{t('setting')}</Link>
                    </Menu.Item>
                </Menu>
      </Sider>
      <Layout>
        <div style={{ position: 'relative' }}>
          <Content>
          <div className='show-main'>
          <div
        className='title'>
              <h1 style={{ fontSize: '32px', color: 'black', textAlign: 'left', lineHeight: '50px', margin: 0, paddingLeft: '24px' }}>
                                    {selectedMenu === 'overview' && <span>{t('overview')}</span>}
                                    {selectedMenu === 'energy' && <span>{t('energy')}</span>}
                                    {selectedMenu === 'map' && <span>{t('map')}</span>}
                                    {selectedMenu === 'log' && <span>{t('log')}</span>}
                                    {selectedMenu === 'history' && <span>{t('history')}</span>}
                                    {selectedMenu === 'home' && <span>{t('home')}</span>}
                                    {selectedMenu === 'media' && <span>{t('media')}</span>}
                                    {selectedMenu === 'todo' && <span>{t('todo')}</span>}
                                    {selectedMenu === 'setting' && <span>{t('setting')}</span>}
                                    {selectedMenu === 'videoMedia' && <span>
                                        <img src={leftarrow}
                                            className='leftarrow'
                                            onClick={handleArrowClick}
                                        />
                                        &nbsp;&nbsp;{t('video')}</span>}
                                </h1>
          </div>
          <div className='show-main'>
              这是能源界面
            </div>
            </div>
          </Content>
        </div>
      </Layout>
    </Layout>
  );
};



