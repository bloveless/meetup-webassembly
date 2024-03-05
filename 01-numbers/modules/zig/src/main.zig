const std = @import("std");

export fn act(a: i32, b: i32) i32 {
    return std.math.pow(i32, a, b);
}

pub fn main() void {}
