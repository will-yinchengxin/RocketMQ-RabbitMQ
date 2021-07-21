<?php
require_once __DIR__ . '/vendor/autoload.php';

use Willyin\Amqp\RabbitMQ;

(new RabbitMQ)->push();
(new RabbitMQ)->cus();
