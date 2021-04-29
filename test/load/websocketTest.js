module.exports = { createTestMessageObject };

async function createTestMessageObject(userContext, events, done) {  
  const randomNumber = Math.floor((Math.random() * 10) + 1);
  const data = {
    body: "Test 12345. " + randomNumber,
    type: "CHAT"
  };
  // set the "data" variable for the virtual user to use in the subsequent action
  userContext.vars.data = data;
  return done();
}

// async function registerAndSetToken(requestParams, context, ee, next) {
//   const registration = await registerChat();
//   console.log('blah')
//   console.log(requestParams.url);
//   return next(); // MUST be called for the scenario to continue
// }

async function registerChat() {
  const options = {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      }
  }

  try {
      const response = await fetch('http://localhost:8080/api/chat/register', options);
      const result = await response.json();
      return result;
  } catch(e) {
      console.error(e);
  }
}