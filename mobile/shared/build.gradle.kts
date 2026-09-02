import com.codingfeline.buildkonfig.compiler.FieldSpec.Type.STRING
import org.jetbrains.kotlin.gradle.dsl.JvmTarget

plugins {
    alias(libs.plugins.kotlinMultiplatform)
    alias(libs.plugins.androidMultiplatformLibrary)
    alias(libs.plugins.buildkonfig)
    alias(libs.plugins.composeMultiplatform)
    alias(libs.plugins.composeCompiler)
}

val envFile = rootDir.parentFile.resolve(".env")

val envValues: Map<String, String> = run {
    check(envFile.exists()) {
        "Arquivo .env não encontrado em ${envFile.path}. Copie .env.example " +
            "para .env na raiz do repositório e preencha os valores antes de " +
            "buildar o mobile."
    }

    val values = mutableMapOf<String, String>()
    envFile.forEachLine { rawLine ->
        val line = rawLine.trim()
        if (line.isEmpty() || line.startsWith("#")) return@forEachLine
        val separatorIndex = line.indexOf('=')
        if (separatorIndex <= 0) return@forEachLine
        values[line.substring(0, separatorIndex).trim()] = line.substring(separatorIndex + 1).trim()
    }
    values
}

fun requiredEnv(key: String): String =
    envValues[key]?.takeIf { it.isNotEmpty() }
        ?: throw GradleException("$key não definida no .env — necessária pro build do mobile")

buildkonfig {
    packageName = "com.squadr"

    defaultConfigs {
        buildConfigField(STRING, "API_BASE_URL", requiredEnv("API_BASE_URL"))
        buildConfigField(STRING, "WS_BASE_URL", requiredEnv("WS_BASE_URL"))
        buildConfigField(STRING, "SUPABASE_URL", requiredEnv("SUPABASE_URL"))
        buildConfigField(STRING, "SUPABASE_PUBLISHABLE_KEY", requiredEnv("SUPABASE_PUBLISHABLE_KEY"))
    }
}

kotlin {
    listOf(
        iosArm64(),
        iosSimulatorArm64()
    ).forEach { iosTarget ->
        iosTarget.binaries.framework {
            baseName = "Shared"
            isStatic = true
        }
    }
    
    android {
       namespace = "com.squadr.shared"
       compileSdk = libs.versions.android.compileSdk.get().toInt()
       minSdk = libs.versions.android.minSdk.get().toInt()
    
       compilerOptions {
           jvmTarget = JvmTarget.JVM_11
       }
       androidResources {
           enable = true
       }
       withHostTest {
           isIncludeAndroidResources = true
       }
       withDeviceTestBuilder {
           sourceSetTreeName = "test"
       }.configure {
           instrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
       }
    }
    
    sourceSets {
        androidMain.dependencies {
            implementation(libs.compose.uiToolingPreview)
            implementation(libs.compose.uiTooling)
            implementation(libs.ktor.client.okhttp)
        }
        iosMain.dependencies {
            implementation(libs.ktor.client.darwin)
        }
        commonMain.dependencies {
            implementation(libs.compose.runtime)
            implementation(libs.compose.foundation)
            implementation(libs.compose.material3)
            implementation(libs.compose.ui)
            implementation(libs.compose.components.resources)
            implementation(libs.compose.uiToolingPreview)
            implementation(libs.androidx.lifecycle.viewmodelCompose)
            implementation(libs.androidx.lifecycle.runtimeCompose)
            implementation(libs.supabase.auth)
            // platform() não existe no KotlinDependencyHandler (o escopo aqui
            // dentro) — só no DependencyHandler padrão do Gradle. Por isso
            // project.dependencies.platform(...), não platform(...) direto.
            // Coordenada em string, não libs.koin.bom: há um bug conhecido do
            // Kotlin (KT-55351) com platform() recebendo valor de version
            // catalog diretamente.
            implementation(project.dependencies.platform("io.insert-koin:koin-bom:${libs.versions.koin.get()}"))
            implementation(libs.koin.core)
        }
        commonTest.dependencies {
            implementation(libs.kotlin.test)
        }
    }
}

dependencies {
    androidRuntimeClasspath(libs.compose.uiTooling)
}