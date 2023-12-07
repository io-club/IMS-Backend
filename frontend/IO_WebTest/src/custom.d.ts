declare module '*.png' {
    const value: any;
    export = value;
  }
  declare module '*.gif' {
    const value: any;
    export = value;
  }

  declare module 'react-leaflet' {
    interface MapContainerProps {
        center: [number, number];
        zoom: number;
        style?: React.CSSProperties; // 为了接受style属性
        children?: React.ReactNode; // 添加children属性
        // 其他属性...
      }
      
      interface TileLayerProps {
        url: string;
        // 其他属性...
      }
      
  
    interface MarkerProps {
      // 根据文档或源代码定义 Marker 的属性
       position: [number, number];
      // 其他属性...
    }
  
    interface PopupProps {
      // 根据文档或源代码定义 Popup 的属性
     position: [number, number];
      // 其他属性...
    }
  
    const MapContainer: React.FC<MapContainerProps>;
    const TileLayer: React.FC<TileLayerProps>;
    const Marker: React.FC<MarkerProps>;
    const Popup: React.FC<PopupProps>;
  }
  
  
  