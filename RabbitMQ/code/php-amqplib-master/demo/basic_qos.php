<?php

include(__DIR__ . '/config.php');

use PhpAmqpLib\Connection\AMQPStreamConnection;

$connection = new AMQPStreamConnection(HOST, PORT, USER, PASS, VHOST);

$channel = $connection->channel();

$channel->queue_declare('qos_queue', false, true, false, false);

//第二个参数代表：每次只消费1条
$channel->basic_qos(null, 1, null);

function process_message($message)
{
    //消费完消息之后进行应答，告诉rabbit我已经消费了，可以发送下一组了
    $message->delivery_info['channel']->basic_ack($message->delivery_info['delivery_tag']);
}

$channel->basic_consume('qos_queue', '', false, false, false, false, 'process_message');

while ($channel->is_consuming()) {
    // After 10 seconds there will be a timeout exception.
    $channel->wait(null, false, 30000);
}