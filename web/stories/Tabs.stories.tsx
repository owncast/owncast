import React from 'react';
import PropTypes from 'prop-types';
import { Tabs, Radio } from 'antd';
import { ComponentStory, ComponentMeta } from '@storybook/react';

const { TabPane } = Tabs;

class TabsExample extends React.Component {
  constructor(props) {
    super(props);

    this.state = { size: 'small' };
  }

  onChange = e => {
    this.setState({ size: e.target.value });
  };

  render() {
    const { size } = this.state;
    const { type } = this.props;

    return (
      <div>
        <Radio.Group value={size} onChange={this.onChange} style={{ marginBottom: 16 }}>
          <Radio.Button value="small">Small</Radio.Button>
          <Radio.Button value="default">Default</Radio.Button>
          <Radio.Button value="large">Large</Radio.Button>
        </Radio.Group>

        <Tabs defaultActiveKey="1" type={type} size={size}>
          <TabPane tab="Card Tab 1" key="1">
            Content of card tab 1
          </TabPane>
          <TabPane tab="Card Tab 2" key="2">
            Content of card tab 2
          </TabPane>
          <TabPane tab="Card Tab 3" key="3">
            Content of card tab 3
          </TabPane>
        </Tabs>
      </div>
    );
  }
}

export default {
  title: 'owncast/Tabs',
  component: Tabs,
} as ComponentMeta<typeof Tabs>;

const Template: ComponentStory<typeof Tabs> = args => <TabsExample {...args} />;

export const Card = Template.bind({});
Card.args = { type: 'card' };

export const Basic = Template.bind({});
Basic.args = { type: '' };

TabsExample.propTypes = {
  type: PropTypes.string,
};

TabsExample.defaultProps = {
  type: '',
};
