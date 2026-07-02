<script setup lang="ts">
import { CommandType } from '@/command-type'
import { getEnv } from '@/env'
import type { Room } from '@/types/room'
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const env = getEnv()
const ws = new WebSocket(new URL('ws', env.apiUrl))
const curUrl = window.location.href
const room = ref<Room>({
  id: '',
  players: [],
})
const route = useRoute()

onMounted(async () => {
  // Connection opened
  ws.addEventListener('open', () => {
    // Join room
    // TODO: Fill player name
    ws.send(
      JSON.stringify({
        type: CommandType.JoinRoom,
        room_id: route.params.room_id as string,
        player_name: '',
      }),
    )
  })

  // Listen for messages
  ws.addEventListener('message', (event) => {
    console.log('Received websocket message:', event.data)
    const roomState = JSON.parse(event.data)
    console.log({ roomState })
    // Update room state
    room.value = {
      id: roomState.id,
      players: roomState.players,
    }
  })

  // Handle errors
  ws.addEventListener('error', (event) => {
    console.error('Webws error:', event)
  })

  // Handle disconnection
  ws.addEventListener('close', (event) => {
    if (event.wasClean) {
      console.log(`Closed cleanly, code=${event.code}, reason=${event.reason}`)
    } else {
      console.log('Connection died')
    }
  })
})

onUnmounted(() => {
  console.log('Closing websocket connection')
  ws.close()
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
              <div class="col-auto d-flex align-items-center">Invite link:</div>
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
