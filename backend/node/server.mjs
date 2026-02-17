import express from 'express';
import { randomBytes } from 'crypto';

const app = express();
app.use(express.json());


const args = {};
for (let i = 2; i < process.argv.length; i++) {
  const arg = process.argv[i];
  if (arg.startsWith('--')) {
    const parts = arg.split('=');
    const key = parts[0].replace(/^--/, '');
    if (parts.length > 1) {
      // Handles --key=value
      args[key] = parts.slice(1).join('=');
    } else {
      // Handles --key (expecting value in next arg)
      // Check if next arg exists and is not another flag
      if (i + 1 < process.argv.length && !process.argv[i + 1].startsWith('--')) {
        args[key] = process.argv[i + 1];
        i++; // Consume the next argument
      } else {
        // Handle boolean flags or flags without value, setting to true if no explicit value
        args[key] = true;
      }
    }
  }
}

const PORT = args.port || 8081;
const FAILURE_RATE = parseInt(args['failure-rate'] || '20', 10);
const SLOW_RATE = parseInt(args['slow-rate'] || '10', 10);

const simulateChaos = (res) => {
  return new Promise((resolve, reject) => {
    const rand = Math.random() * 100;
    if (rand < SLOW_RATE) {
      setTimeout(resolve, 6000);
    } else if (rand < SLOW_RATE + FAILURE_RATE) {
      res.status(500).json({ error: 'Upstream service unavailable' });
      reject('Chaos Failure');
    } else {
      resolve();
    }
  });
};


const availableSKUs = [
  { id: "C1-R1GB-D20GB", sku: "C1-R1GB-D20GB", type: "virtual-machine", name: "Micro", cpu: 1, ram: 1, disk: 20, price_hourly: 0.27, price_monthly: 180 },
  { id: "C2-R4GB-D80GB", sku: "C2-R4GB-D80GB", type: "virtual-machine", name: "Standard", cpu: 2, ram: 4, disk: 80, price_hourly: 1.1, price_monthly: 750 },
  { id: "C4-R8GB-D160GB", sku: "C4-R8GB-D160GB", type: "virtual-machine", name: "Performance", cpu: 4, ram: 8, disk: 160, price_hourly: 2.2, price_monthly: 1500 },
  { id: "C8-R32GB-D320GB", sku: "C8-R32GB-D320GB", type: "virtual-machine", name: "Pro Max", cpu: 8, ram: 32, disk: 320, price_hourly: 5.2, price_monthly: 3500 },
  { id: "C8-R16GB-D512GB", sku: "C8-R16GB-D512GB", type: "dedicated", name: "Metal Alpha", cpu: 8, ram: 16, disk: 512, price_hourly: 18, price_monthly: 12000 },
  { id: "C16-R64GB-D1024GB", sku: "C16-R64GB-D1024GB", type: "dedicated", name: "Metal Beta", cpu: 16, ram: 64, disk: 1024, price_hourly: 42, price_monthly: 28000 },
  { id: "C32-R128GB-D2048GB", sku: "C32-R128GB-D2048GB", type: "dedicated", name: "Metal Gamma", cpu: 32, ram: 128, disk: 2048, price_hourly: 90, price_monthly: 60000 },
  { id: "C64-R256GB-D4096GB", sku: "C64-R256GB-D4096GB", type: "dedicated", name: "Metal Omega", cpu: 64, ram: 256, disk: 4096, price_hourly: 180, price_monthly: 120000 },
];

const resources = {};


app.get('/v1/skus', (req, res) => {
  res.json({ skus: availableSKUs });
});

app.post('/v1/availability', (req, res) => {
  const { sku } = req.body;
  if (!sku) return res.status(400).json({ error: 'Invalid JSON' });

  let available = availableSKUs.some(s => s.sku === sku);
  if (sku === 'OUT-OF-STOCK') available = false;

  res.json({ sku, available });
});

app.get('/v1/resources', (req, res) => {
  res.json({ resources: Object.values(resources) });
});

app.post('/v1/resources', async (req, res) => {
  try {
    await simulateChaos(res);
  } catch { return; }

  const { sku } = req.body;
  const isValidSKU = availableSKUs.some(s => s.sku === sku) || sku === 'test';

  if (!isValidSKU) {
    return res.status(400).json({ error: 'Invalid SKU' });
  }

  setTimeout(() => {
    const id = `i-${randomBytes(4).toString('hex')}`;
    const resource = {
      id,
      sku,
      status: 'running',
      ip: `10.0.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`,
      created_at: new Date().toISOString()
    };
    resources[id] = resource;
    res.json(resource);
  }, 500 + Math.random() * 1000);
});

app.get('/v1/resources/:id', (req, res) => {
  const resource = resources[req.params.id];
  if (!resource) return res.status(404).json({ error: 'Resource not found' });
  res.json(resource);
});

app.post('/v1/resources/:id/power', async (req, res) => {
  const resource = resources[req.params.id];
  if (!resource) return res.status(404).json({ error: 'Resource not found' });

  try {
    await simulateChaos(res);
  } catch { return; }

  const { action } = req.body;
  if (!['on', 'off'].includes(action)) {
    return res.status(400).json({ error: 'Invalid action' });
  }

  resource.status = action === 'on' ? 'running' : 'stopped';

  res.json({ status: 'success', state: action });
});

app.listen(PORT, () => {
  console.log(`🔌 Infra Service running on :${PORT}`);
});
