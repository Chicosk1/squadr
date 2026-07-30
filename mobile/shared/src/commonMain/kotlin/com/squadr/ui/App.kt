package com.squadr.ui

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.safeContentPadding
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier

/**
 * Placeholder de verificação da Fase 2 — confirma que o build e o runtime
 * (Android + iOS) funcionam de ponta a ponta. Substituído pela navegação real
 * no Bloco 0 da Fase 4 (login → onboarding → abas → chat).
 */
@Composable
fun App() {
    MaterialTheme {
        Box(
            modifier = Modifier.safeContentPadding().fillMaxSize(),
            contentAlignment = Alignment.Center,
        ) {
            Text("Squadr")
        }
    }
}
