use extism_pdk::*;

#[host_fn]
extern "ExtismHost" {
    fn owncast_send_chat(text: String);
}

// A self-contained module has no SDK to bake the sidecar into register(). It
// must read the exact manifest Owncast loaded from the reserved Extism config.
#[plugin_fn]
pub fn register() -> FnResult<String> {
    Ok(config::get("manifest")?
        .ok_or_else(|| Error::msg("Owncast did not inject the plugin manifest"))?)
}

// Echo the raw event envelope. The integration test looks for its unique chat
// probe in the bot response, proving Owncast loaded the native module, installed
// the subscription returned by register(), dispatched the event, and allowed
// the module to call an Owncast host function.
#[plugin_fn]
pub fn on_event(envelope: String) -> FnResult<()> {
    unsafe { owncast_send_chat(envelope)? };
    Ok(())
}
