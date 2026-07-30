package com.squadr.platform

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform