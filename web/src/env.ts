import type { Env } from './global'

function getEnv(): Env {
  if (!window.__ENV) {
    console.log('window.__ENV is missing. Use default env')
    return { env: 'local', apiUrl: 'http://localhost:3000/api/' }
  }

  return window.__ENV
}

export { getEnv }
