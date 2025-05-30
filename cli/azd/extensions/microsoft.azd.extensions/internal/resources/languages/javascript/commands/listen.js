const { Command } = require('commander');
const AzdClient = require('../azdClient');
const { EventManager } = require('../eventManager');
const logger = require('../logger');

function createListenCommand() {
  const cmd = new Command('listen');
  cmd.description('Starts the extension and listens for events.');

  cmd.action(async () => {
    const client = new AzdClient();
    const eventManager = new EventManager(client);

    const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));

    await eventManager.addProjectEventHandler('preprovision', async () => {
      for (let i = 1; i <= 20; i++) {
        console.log(`[preprovision] Doing important work... step ${i}`);
        await sleep(200);
      }
    });

    await eventManager.addServiceEventHandler('prepackage', async () => {
      for (let i = 1; i <= 20; i++) {
        console.log(`[prepackage] Doing important work... step ${i}`);
        await sleep(200);
      }
    });

    try {
      await eventManager.receive();
    } catch (err) {
      logger.error('Error while receiving events:', { error: err.message, stack: err.stack });
    }
  });

  return cmd;
}

module.exports = { createListenCommand };
