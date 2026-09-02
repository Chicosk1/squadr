import SwiftUI
import Shared

@main
struct iOSApp: App {
    init() {
        KoinKt.initKoin()
    }

    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}