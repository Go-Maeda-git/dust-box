#!/usr/bin/env python3
import rclpy
from rclpy.node import Node
from sensor_msgs.msg import LaserScan
import math

class LidarSubscriber(Node):

    def __init__(self):
        super().__init__('lidar_subscriber_node')
        self.subscription = self.create_subscription(
            LaserScan,
            '/scan', 
            self.listener_callback,
            10)
        self.get_logger().info('LiDAR Subscriber Node has been started.')

    def listener_callback(self, msg):
        min_distance = float('inf')
        min_index = -1
        valid_ranges = []

        for i, distance in enumerate(msg.ranges):
            if not math.isinf(distance) and not math.isnan(distance) and msg.range_min < distance < msg.range_max:
                valid_ranges.append(distance)
                if distance < min_distance:
                    min_distance = distance
                    min_index = i

        if min_index != -1:
            angle = msg.angle_min + min_index * msg.angle_increment
            angle_degrees = math.degrees(angle)
            self.get_logger().info(
                f'Closest obstacle: {min_distance:.2f} [m] at index {min_index} ({angle_degrees:.2f} [deg])'
            )
        else:
            self.get_logger().info('No valid obstacles detected within range.')

def main(args=None):
    rclpy.init(args=args)
    lidar_subscriber = LidarSubscriber()
    rclpy.spin(lidar_subscriber)
    lidar_subscriber.destroy_node()
    rclpy.shutdown()

if __name__ == '__main__':
    main()
