// add more to the promises later.
class Config {
  async init() {
    const configFileLocation = "https://goth.land/config";
    
    try {
      const response = await fetch(configFileLocation);
      const configData = await response.json();
      Object.assign(this, configData);
      return this;
    } catch(error) {
      console.log(error);
      // No config file present.  That's ok.  It's not required.
    }
  }
}