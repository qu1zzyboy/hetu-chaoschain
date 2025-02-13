import { Flex, Image } from 'antd';
import React from 'react';
import {  MinusCircleOutlined } from '@ant-design/icons';
import image16 from '@/assets/image16.png'
import image14 from '@/assets/image14.png'
interface Props {
  status?:number;
}

// 脚手架示例组件
const Status: React.FC<Props> = (props) => {
  const { status } = props;
  const handleStatus = () => {
    let color,text = '';
    switch (status) {
      case 1:
      case 4:{
        color = '#FDBF0780';
        text = 'Failed';
      }
      break;
      case 2: {
        color = '#4DCA8280';
        text = 'Processing';
      }
      break;
      case 3: {
        color = '#D2EDFE';
        text = 'Passed';
      }
      break;
    }
    return { color,text}
  };
  return (
    <>
      {status && <Flex justify={'space-between'} align={'center'} style={{padding:'10px',background:handleStatus().color,borderRadius:'14px',height:'28px',fontSize:'15px',color:'#777575'}} >
        {(status === 1 || status === 4) && <MinusCircleOutlined style={{ marginRight: '4px' }} />}
        {status === 2 && <Image preview={false} style={{marginTop:'-4px',paddingRight:'4px'}} src={image14}/>}
        {status === 3 && <Image preview={false} style={{marginTop:'-4px',paddingRight:'4px'}} src={image16}/>}
        {handleStatus().text}
      </Flex>}
    </>
  );
};

export default Status;