import { Flex } from 'antd';
import React, { useEffect } from 'react';
import { CheckCircleFilled, CloseCircleFilled } from '@ant-design/icons';
import './index.less';

interface Props {
  status?: number;
}

interface ItemClass {
  text: string;
  icon?: React.ReactNode | null;
  textColor: string;
  progressColor: string;
}

const BottomProgressBar: React.FC<Props> = (props) => {
  const { status } = props;
  const [itemClasses, setItemClasses] = React.useState<ItemClass[]>([]);
  useEffect(() => {
    //忽略（状态码1）、讨论中（状态码2）、接受（状态码3）、拒绝（状态码4）
    setItemClasses([
      {
        text:'Draft',
        icon:status === 1 ?<CloseCircleFilled />: <CheckCircleFilled />,
        textColor: status === 1 ?'text-red': 'text-green',
        progressColor: status === 1 ?'progress-red': 'progress-green-1',
      },
      {
        text:'Discussion',
        icon: (status === 3 || status === 4) ? <CheckCircleFilled />:null,
        textColor: status === 1?'text-gray':'text-green',
        progressColor: status === 1?'progress-gray':(status === 2 || status === 4?'progress-green-3':'progress-green-3'),
      },
      {
        text:'Decision',
        icon: (status === 3 && <CheckCircleFilled />) || (status === 4 && <CloseCircleFilled />),
        textColor: status === 3 ?'text-green':(status === 4 ?'text-red':'text-gray'),
        progressColor: status === 3 ?'progress-green-3':(status === 4 ?'progress-red':'progress-gray'),
      }]);
  }, [status]);

  return (
    <Flex justify={'space-between'} vertical align={'center'} style={{ width: `100%` }}>
      <Flex className={'progress-item-base'} style={{ width: '100%' }}>
        {itemClasses?.map((item, index) => <Flex key={`progress-name-${index}`} className={`${item?.textColor}`} justify={'center'} flex={1}>
          <span>{`${item?.text}  `}{item?.icon}</span>
        </Flex>)}
      </Flex>
      <Flex className='progress-bar'>
        {itemClasses?.map((item, index) => <div key={`progress-${index}`} className={`progress-segment ${item?.progressColor}`}></div>)}
      </Flex>
    </Flex>
  );
};

export default BottomProgressBar;