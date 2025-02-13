import React, { useEffect } from 'react';
import { Divider, Flex, Image, Typography } from 'antd';
import EChartsReact from '@/components/BaseCharts';
import { EChartsOption } from 'echarts';
import './style.less'

interface Props {
  status: number;
  showLine: boolean;
  decisionStep:AGENTAPI.DecisionStep
}

const avatarContext = require.context(
  '@/assets/head', // 目录路径
  false, // 不递归子目录
  /\.(png|jpe?g|svg)$/, // 匹配格式
);
const avatarList = avatarContext.keys().map(key => avatarContext(key));

const Discussions: React.FC<Props> = (props) => {
  const { status,showLine = false,decisionStep } = props;
  const [passPercent, setPassPercent] = React.useState(0);
  useEffect(() => {
    if (decisionStep?.decisionPass > 0 || decisionStep?.decisionReject > 0) {
      setPassPercent(Math.ceil(decisionStep?.decisionPass / (decisionStep?.decisionPass + decisionStep?.decisionReject))*100);
    }

  }, []);

  const itemCiaOption: EChartsOption = {
    tooltip: {
      trigger: 'item',
      formatter:(obj ) => {
        return `${'name' in obj ? obj?.name :''}: ${'value' in obj ? obj?.value :''}%`
      }
    },
    series: [
      {
        type: 'pie',
        radius: '70%',
        itemStyle: {
          color: function (params: { dataIndex: number }) {
            const colorList: string[] = ['#5798F7','#E36681' ];
            return colorList[params.dataIndex];
          },
        },
        data: [
          { value: passPercent, name: `Pass` },
          { value: 100-passPercent, name: 'Reject' },
        ],
        label: {
          formatter: (obj) => {
            return `${obj?.name} ${obj?.value}%`
          }
        },
      }
    ]
  };

  return (
    <div>
      <Flex className={'discussions-title'}>Discussions</Flex>
      {decisionStep?.discussions?.map((item, index) => <Flex className={'discussions-item'} vertical key={`discussions-${index}`}>
        <Flex align={'center'} className={'discussions-item-header'}>
          <Image width={80} height={80} preview={false}
                 src={avatarList[index % 4]} />
          <span style={{marginLeft:'8px'}}>{item?.speaker_name}</span>
        </Flex>
        <Typography.Paragraph className={'discussions-item-content'}>
          {item?.data}
        </Typography.Paragraph>
      </Flex>)}
      {(decisionStep?.decisionPass > 0 || decisionStep?.decisionReject > 0) && <>
        <Flex className={'discussions-title'}>Decision</Flex>
        <Flex>
          <EChartsReact
            option={itemCiaOption}
            style={{ height: 238, width: '380px' }}
          />
        </Flex>
      </>}
      {showLine && <Divider style={{background:'#000000'}} />}
    </div>
  );
};

export default Discussions;