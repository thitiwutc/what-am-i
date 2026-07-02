<script setup lang="ts">
import { getEnv } from '@/env'
import type { Room } from '@/types/room'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const curUrl = window.location.href
const room = ref<Room>({
  id: '',
  players: [],
})
const route = useRoute()

onMounted(async () => {
  const env = getEnv()
  const url = new URL('rooms/' + route.params.room_id + '/players', env.apiUrl)
  const resp = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      player_name: '',
    }),
  })

  const respBody = await resp.json()

  room.value.id = route.params.room_id as string
  room.value.players = respBody.data.players
})
</script>

<template>
  <div class="container d-flex justify-content-center align-items-center vh-100">
    <div class="row w-100 justify-content-center">
      <div
        class="col-12 col-sm-10 col-md-8 col-lg-6 col-xl-4 shadow p-3 rounded border"
        style="background: white"
      >
        <div class="row gy-2">
          <h2 class="col-12">Room ID: {{ room.id }}</h2>
          <div class="col-12">
            <div class="row">
              <div class="col-auto">Invite link:</div>
              <div class="col">
                <input class="form-control" type="text" readonly :value="curUrl" />
              </div>
            </div>
          </div>
          <div class="col-12">
            <div class="row">
              <div class="col-12">
                <b>{{ room.players.length > 1 ? 'Players' : 'Players' }}</b>
              </div>
              <div class="col-12" v-for="player in room.players" :key="player.id">
                {{ player.name }}
              </div>
              <div class="div d-flex justify-content-end">
                <button class="btn btn-primary">Start</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
