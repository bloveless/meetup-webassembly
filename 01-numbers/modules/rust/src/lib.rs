#[cfg_attr(all(target_arch = "wasm32"), export_name = "act")]
#[no_mangle]
pub extern "C" fn act(x: i64, y: i64) -> i64 {
    x * y
}
