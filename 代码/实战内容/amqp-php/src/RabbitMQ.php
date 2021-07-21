<?php
/**
 * AMQP
 *
 * @package   Willyin\Amqp
 * @author    Will  <826895143@qq.com>
 * @copyright Copyright (C) 2021 Will
 */

namespace Willyin\Amqp;


class RabbitMQ
{
    protected $connect;

    protected $config = [
        'host' => '127.0.0.1',
        'vhost' => '/',
        'port' => 5672,
        'login' => 'guest',
        'password' => 'guest'
    ];

    /**
     * 发送消息(路由模式模式)
     *
     * 每个消费者监听自己的队列，并且设置带统配符的 routingkey
     * 生产者将消息发给broker，由交换机根据 routingkey 来转发消息到指定的队列
     */
    public function push()
    {
        $cnn = new \AMQPConnection($this->config);
        if (!$cnn->connect()) {
            echo "Can't connect to the test";
            exit();
        }
        $ch = new \AMQPChannel($cnn);
        $ex = new \AMQPExchange($ch);

        //消息的路由键，一定要和消费者端一致
        $routingKeyOne = 'key_1';
        $routingKeyTwo = 'key_2';

        //交换机名称，一定要和消费者端一致，
        $exchangeName = 'exchange_1';

        $ex->setName($exchangeName);
        $ex->setType(AMQP_EX_TYPE_DIRECT);
        $ex->setFlags(AMQP_DURABLE);
        $ex->declareExchange();

        //创建一个消息队列
        $q = new \AMQPQueue($ch);
        //设置队列名称
        $q->setName('queue_1');
        //设置队列持久
        $q->setFlags(AMQP_DURABLE);
        //声明消息队列
        $q->declareQueue();
        //交换机和队列通过$routingKey进行绑定
        $q->bind($ex->getName(), $routingKeyOne);

        /** 创建了多个queue,共享交换机中的信息
        $qOne = new \AMQPQueue($ch);
        //设置队列名称
        $qOne->setName('queue_2');
        //设置队列持久
        $qOne->setFlags(AMQP_DURABLE);
        //声明消息队列
        $qOne->declareQueue();
        //交换机和队列通过$routingKey进行绑定
        $qOne->bind($ex->getName(), $routingKeyTwo);
        **/

        //创建10个消息
        for ($i = 1; $i <= 10; $i++) {
            //消息内容
            $msg = array(
                'data' => 'message_' . $i,
                'hello' => 'world',
            );
            //发送消息到交换机，并返回发送结果
            //delivery_mode:2声明消息持久，持久的队列+持久的消息在RabbitMQ重启后才不会丢失
            echo "Send Message:" . $ex->publish(json_encode($msg), $routingKeyOne
                    , AMQP_NOPARAM, array('delivery_mode' => 2)) . "\n";
            //代码执行完毕后进程会自动退出
        }
        /**
         * 声明队列并声明交换机 -> 创建连接 -> 创建通道 -> 通道声明交换机 -> 通道声明队列 -> 通过通道使队列绑定到交换机并指定该队列的routingkey（通配符）
         *  -> 制定消息 -> 发送消息并指定routingkey（通配符）
         * */
    }

    /**
     * 消费消息
     */
    public function cus()
    {
        //连接
        $cnn = new \AMQPConnection($this->config);
        if (!$cnn->connect()) {
            echo "Cannot connect to the broker";
            exit();
        }

        //在连接内创建一个通道
        $ch = new \AMQPChannel($cnn);

        //创建一个交换机
        $ex = new \AMQPExchange($ch);

        //声明路由键
        $routingKey = 'key_1';

        //声明交换机名称
        $exchangeName = 'exchange_1';

        //设置交换机名称
        $ex->setName($exchangeName);

        //设置交换机类型
        //AMQP_EX_TYPE_DIRECT:直连交换机
        //AMQP_EX_TYPE_FANOUT:扇形交换机
        //AMQP_EX_TYPE_HEADERS:头交换机
        //AMQP_EX_TYPE_TOPIC:主题交换机
        $ex->setType(AMQP_EX_TYPE_DIRECT);

        //设置交换机持
        $ex->setFlags(AMQP_DURABLE);

        //声明交换机
        $ex->declareExchange();

        //创建一个消息队列
        $q = new \AMQPQueue($ch);

        //设置队列名称
        $q->setName('queue_1');

        //设置队列持久
        $q->setFlags(AMQP_DURABLE);

        //声明消息队列
        $q->declareQueue();

        //交换机和队列通过$routingKey进行绑定
        $q->bind($ex->getName(), $routingKey);

        //设置消息队列消费者回调方法，
        $q->consume(function ($envelope, $queue) {
            //休眠两秒，
            sleep(1);
            //echo消息内容
            echo $envelope->getBody() . "\n";
            //显式确认，队列收到消费者显式确认后，会删除该消息
            $queue->ack($envelope->getDeliveryTag());
        });
    }
}