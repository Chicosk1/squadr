package com.squadr

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.runtime.Composable
import androidx.compose.ui.tooling.preview.Preview
import com.squadr.di.handleAuthDeeplink
import com.squadr.di.initKoin
import com.squadr.ui.App

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge()
        super.onCreate(savedInstanceState)
        initKoin()
        handleAuthDeeplink(intent)

        setContent {
            App()
        }
    }

    // launchMode="singleTask" (AndroidManifest) faz o retorno do OAuth chegar
    // aqui, não num onCreate novo — sem isso o deep link seria perdido se a
    // Activity já estivesse aberta quando o navegador voltar.
    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        setIntent(intent)
        handleAuthDeeplink(intent)
    }
}

@Preview
@Composable
fun AppAndroidPreview() {
    App()
}