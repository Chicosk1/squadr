package com.squadr.di

import android.content.Intent
import io.github.jan.supabase.SupabaseClient
import io.github.jan.supabase.auth.handleDeeplinks
import org.koin.mp.KoinPlatform

// Wrapper pra que MainActivity (módulo androidApp) não precise depender
// direto de Koin nem do auth-kt — só desse tipo Android comum (Intent).
fun handleAuthDeeplink(intent: Intent) {
    KoinPlatform.getKoin().get<SupabaseClient>().handleDeeplinks(intent)
}
