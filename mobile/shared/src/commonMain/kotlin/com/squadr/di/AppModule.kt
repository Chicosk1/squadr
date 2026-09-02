package com.squadr.di

import com.squadr.BuildKonfig
import io.github.jan.supabase.SupabaseClient
import io.github.jan.supabase.auth.Auth
import io.github.jan.supabase.createSupabaseClient
import org.koin.dsl.module

// scheme/host precisam ser idênticos ao AndroidManifest.xml, ao Info.plist
// do iosApp e ao cadastrado no painel do Supabase — ver ADR-010.
private const val OAUTH_CALLBACK_SCHEME = "com.squadr"
private const val OAUTH_CALLBACK_HOST = "login-callback"

val appModule = module {
    single<SupabaseClient> {
        createSupabaseClient(
            supabaseUrl = BuildKonfig.SUPABASE_URL,
            supabaseKey = BuildKonfig.SUPABASE_PUBLISHABLE_KEY,
        ) {
            install(Auth) {
                scheme = OAUTH_CALLBACK_SCHEME
                host = OAUTH_CALLBACK_HOST
            }
        }
    }
}
