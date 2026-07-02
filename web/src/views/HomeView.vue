<script setup lang="ts">
import { getEnv } from '@/env'
import { useRouter } from 'vue-router'

const router = useRouter()

async function createRoom() {
  const env = getEnv()
  const url = new URL('rooms/', env.apiUrl)
  const resp = await fetch(url, {
    method: 'POST',
  })
  const respBody = await resp.json()
  const roomId = respBody.data.room_id

  router.push('/room/' + roomId)
}
</script>

<template>
  <div class="container d-flex justify-content-center align-items-center vh-100">
    <div class="row w-100 justify-content-center">
      <div
        class="col-12 col-sm-10 col-md-8 col-lg-6 col-xl-4 shadow py-2 rounded border"
        style="background: white"
      >
        <div class="row">
          <div class="col-12 pt-1">
            <input
              type="text"
              class="form-control"
              placeholder="Player name"
              aria-label="Player name"
            />
          </div>
        </div>
        <div class="row mt-2 gy-2">
          <div class="col-12 text-center">
            <button class="btn btn-primary" @click="createRoom">Create room</button>
          </div>
          <div class="col-12 text-center">
            <button class="btn btn-secondary">Join room</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
