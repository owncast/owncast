class Config {

  constructor() {
    this.init();
  }

  async init() {
    const configFileLocation = "js/config.json";
    
    try {
      const response = await fetch(configFileLocation);
      const configData = await response.json();
      Object.assign(this, configData);
    } catch(error) {
      console.log(error);
      // No config file present.  That's ok.  It's not required.
    }
  }
}