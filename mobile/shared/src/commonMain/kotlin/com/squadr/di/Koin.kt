package com.squadr.di

import org.koin.core.context.startKoin

private var koinStarted = false

// Chamada uma vez, no início do processo (MainActivity.onCreate no Android,
// iOSApp.init no iOS). O guard evita KoinAppAlreadyStartedException se o
// ponto de entrada da plataforma rodar mais de uma vez.
fun initKoin() {
    if (koinStarted) return
    koinStarted = true

    startKoin {
        modules(appModule)
    }
}
